from __future__ import annotations

from pathlib import Path
from typing import Literal
import json
import os

from crewai import Agent, Crew, LLM, Process, Task
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
import yaml


load_dotenv()

BASE_DIR = Path(__file__).resolve().parent
CONFIG_DIR = BASE_DIR / "config"

app = FastAPI(title="selfCook AI Service")


class ChatMessage(BaseModel):
    role: Literal["user", "assistant", "system"] = "user"
    content: str


class MenuItem(BaseModel):
    groupId: int
    groupTitle: str
    cutoffAt: str
    fulfillmentMode: str
    groupNotice: str = ""
    productId: int
    groupItemId: int
    name: str
    skuName: str
    price: float
    originalPrice: float = 0
    stockAvailable: int
    limitPerOrder: int = 0
    categoryName: str = ""
    description: str = ""


class RecommendRequest(BaseModel):
    message: str = Field(min_length=1, max_length=500)
    history: list[ChatMessage] = Field(default_factory=list)
    menuItems: list[MenuItem] = Field(default_factory=list)


class RecommendResponse(BaseModel):
    reply: str


def load_yaml(name: str) -> dict:
    with open(CONFIG_DIR / name, encoding="utf-8") as file:
        return yaml.safe_load(file)


def deepseek_llm() -> LLM:
    api_key = os.getenv("DEEPSEEK_API_KEY", "")
    if not api_key:
        raise HTTPException(status_code=500, detail="DEEPSEEK_API_KEY is not configured")

    return LLM(
        model=os.getenv("DEEPSEEK_MODEL", "deepseek/deepseek-chat"),
        api_key=api_key,
        base_url=os.getenv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),
        temperature=0.3,
    )


def build_history_text(history: list[ChatMessage]) -> str:
    if not history:
        return "无"

    recent = history[-8:]
    lines = []
    for item in recent:
        role = "用户" if item.role == "user" else "助手"
        content = item.content.strip()
        if content:
            lines.append(f"{role}: {content}")
    return "\n".join(lines) or "无"


def build_menu_text(menu_items: list[MenuItem]) -> str:
    if not menu_items:
        return "当前没有可售菜单。"

    lines = []
    for item in menu_items:
        limit = f"，单次限购{item.limitPerOrder}份" if item.limitPerOrder else ""
        desc = f"，说明：{item.description}" if item.description else ""
        lines.append(
            "- "
            f"团购[{item.groupTitle}] "
            f"菜品ID:{item.groupItemId} "
            f"{item.name}/{item.skuName}，"
            f"分类:{item.categoryName or '未分类'}，"
            f"价格:{item.price:.2f}元，"
            f"库存:{item.stockAvailable}份{limit}，"
            f"履约:{item.fulfillmentMode}，"
            f"截单:{item.cutoffAt}"
            f"{desc}"
        )
    return "\n".join(lines)


def build_menu_json(menu_items: list[MenuItem]) -> str:
    return json.dumps([item.model_dump() for item in menu_items], ensure_ascii=False)


def build_menu_crew() -> Crew:
    agents_config = load_yaml("agents.yaml")
    tasks_config = load_yaml("tasks.yaml")
    llm = deepseek_llm()

    menu_planner = Agent(
        config=agents_config["menu_planner"],
        llm=llm,
        verbose=False,
    )
    recommendation_task = Task(
        config=tasks_config["recommend_menu_task"],
        agent=menu_planner,
    )

    return Crew(
        agents=[menu_planner],
        tasks=[recommendation_task],
        process=Process.sequential,
        verbose=False,
    )


@app.get("/health")
def health() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/crew/recommend", response_model=RecommendResponse)
def recommend(req: RecommendRequest) -> RecommendResponse:
    crew = build_menu_crew()
    result = crew.kickoff(
        inputs={
            "message": req.message.strip(),
            "history": build_history_text(req.history),
            "menu_items": build_menu_text(req.menuItems),
            "menu_items_json": build_menu_json(req.menuItems),
        }
    )
    return RecommendResponse(reply=str(result))
