const API_BASE = "/api/v1";

const state = {
  mode: "login",
  token: localStorage.getItem("fz_token") || "",
  userId: localStorage.getItem("fz_user_id") || "",
  username: localStorage.getItem("fz_username") || "",
  items: [],
  topCursor: "",
  bottomCursor: "",
  hasMore: false,
  feedLoading: false,
  menuContentId: "",
  wheelAt: 0,
};

const $ = (id) => document.getElementById(id);

const el = {
  accountStrip: $("accountStrip"),
  sessionStrip: $("sessionStrip"),
  composer: $("composer"),
  loginTab: $("loginTab"),
  registerTab: $("registerTab"),
  authForm: $("authForm"),
  authSubmit: $("authSubmit"),
  usernameInput: $("usernameInput"),
  passwordInput: $("passwordInput"),
  logoutButton: $("logoutButton"),
  avatar: $("avatar"),
  sessionName: $("sessionName"),
  sessionId: $("sessionId"),
  postText: $("postText"),
  postCounter: $("postCounter"),
  publishButton: $("publishButton"),
  refreshButton: $("refreshButton"),
  cameraButton: $("cameraButton"),
  notifyButton: $("notifyButton"),
  loadOlderButton: $("loadOlderButton"),
  feedList: $("feedList"),
  stream: $("stream"),
  refreshHint: $("refreshHint"),
  followeeInput: $("followeeInput"),
  followButton: $("followButton"),
  unfollowButton: $("unfollowButton"),
  apiStatus: $("apiStatus"),
  cursorStatus: $("cursorStatus"),
  lastAction: $("lastAction"),
  toast: $("toast"),
  postMenu: $("postMenu"),
  deletePostButton: $("deletePostButton"),
};

function mountIcons() {
  const icons = {
    bell: '<path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"></path><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"></path>',
    camera: '<path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"></path><circle cx="12" cy="13" r="3"></circle>',
    "rotate-cw": '<path d="M21 12a9 9 0 1 1-2.64-6.36"></path><path d="M21 3v6h-6"></path>',
    "more-horizontal": '<circle cx="12" cy="12" r="1"></circle><circle cx="19" cy="12" r="1"></circle><circle cx="5" cy="12" r="1"></circle>',
    "trash-2": '<path d="M3 6h18"></path><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"></path><path d="M10 11v6"></path><path d="M14 11v6"></path>',
    heart: '<path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.29 1.51 4.04 3 5.5l7 7Z"></path>',
    play: '<polygon points="6 3 20 12 6 21 6 3"></polygon>',
  };
  document.querySelectorAll("i[data-lucide]").forEach((node) => {
    const name = node.getAttribute("data-lucide");
    if (!icons[name]) return;
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.setAttribute("viewBox", "0 0 24 24");
    svg.setAttribute("fill", "none");
    svg.setAttribute("stroke", "currentColor");
    svg.setAttribute("stroke-width", "2");
    svg.setAttribute("stroke-linecap", "round");
    svg.setAttribute("stroke-linejoin", "round");
    svg.setAttribute("aria-hidden", "true");
    svg.innerHTML = icons[name];
    node.replaceWith(svg);
  });
}

function setMode(mode) {
  state.mode = mode;
  el.loginTab.classList.toggle("active", mode === "login");
  el.registerTab.classList.toggle("active", mode === "register");
  el.authSubmit.textContent = mode === "login" ? "登录" : "注册";
  el.passwordInput.autocomplete = mode === "login" ? "current-password" : "new-password";
}

function setBusy(target, busy) {
  if (target) target.disabled = busy;
}

function showToast(message) {
  el.toast.textContent = message;
  el.toast.classList.remove("hidden");
  window.clearTimeout(showToast.timer);
  showToast.timer = window.setTimeout(() => el.toast.classList.add("hidden"), 2400);
}

function setLastAction(message) {
  el.lastAction.textContent = message;
}

function setHint(message, active = false) {
  el.refreshHint.textContent = message;
  el.refreshHint.classList.toggle("active", active);
}

async function request(path, options = {}) {
  const headers = {
    "Content-Type": "application/json",
    ...(options.headers || {}),
  };
  if (state.token) {
    headers.Authorization = `Bearer ${state.token}`;
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  const text = await response.text();
  const body = text ? JSON.parse(text) : null;
  if (!response.ok) {
    throw new Error(body?.message || `HTTP ${response.status}`);
  }
  return body?.data;
}

async function checkHealth() {
  try {
    const response = await fetch("/healthz");
    el.apiStatus.textContent = response.ok ? "API 正常" : `API ${response.status}`;
  } catch (error) {
    el.apiStatus.textContent = "API 不可达";
  }
}

function updateSessionView() {
  const signedIn = Boolean(state.token && state.userId);
  el.accountStrip.classList.toggle("hidden", signedIn);
  el.sessionStrip.classList.toggle("hidden", !signedIn);
  el.composer.classList.toggle("hidden", !signedIn);
  el.loadOlderButton.disabled = !signedIn;
  el.refreshButton.disabled = !signedIn;
  el.cameraButton.disabled = !signedIn;

  el.sessionName.textContent = state.username || "已登录";
  el.sessionId.textContent = state.userId ? `ID ${state.userId}` : "-";
  el.avatar.textContent = (state.username || "U").slice(0, 1).toUpperCase();
}

function saveSession(data, username) {
  state.token = data.token;
  state.userId = String(data.user_id);
  state.username = username;
  localStorage.setItem("fz_token", state.token);
  localStorage.setItem("fz_user_id", state.userId);
  localStorage.setItem("fz_username", state.username);
  updateSessionView();
}

function clearSession() {
  state.token = "";
  state.userId = "";
  state.username = "";
  state.items = [];
  state.topCursor = "";
  state.bottomCursor = "";
  state.hasMore = false;
  localStorage.removeItem("fz_token");
  localStorage.removeItem("fz_user_id");
  localStorage.removeItem("fz_username");
  updateSessionView();
  renderFeed();
  updateCursorStatus();
  setHint("准备就绪");
}

async function submitAuth(event) {
  event.preventDefault();
  const username = el.usernameInput.value.trim();
  const password = el.passwordInput.value;
  if (!username || !password) return;

  setBusy(el.authSubmit, true);
  try {
    const data = await request(state.mode === "login" ? "/auth/login" : "/auth/register", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    });
    saveSession(data, username);
    el.passwordInput.value = "";
    setLastAction("登录成功");
    showToast(state.mode === "login" ? "登录成功" : "注册成功");
    await loadFeed("latest", { quiet: true });
  } catch (error) {
    showToast(error.message);
    setLastAction("登录失败");
  } finally {
    setBusy(el.authSubmit, false);
  }
}

async function publishPost() {
  const content = el.postText.value.trim();
  if (!content) {
    showToast("请输入内容");
    return;
  }

  setBusy(el.publishButton, true);
  try {
    const data = await request("/posts", {
      method: "POST",
      body: JSON.stringify({ content_text: content }),
    });
    el.postText.value = "";
    updateCounter();
    showToast(`已发表 ${data.content_id}`);
    setLastAction("发表动态");
    setHint("向上滚轮或点刷新可查看最新", true);
  } catch (error) {
    showToast(error.message);
    setLastAction("发表失败");
  } finally {
    setBusy(el.publishButton, false);
  }
}

async function followUser(shouldFollow) {
  const followeeId = el.followeeInput.value.trim();
  if (!/^\d+$/.test(followeeId)) {
    showToast("请输入用户 ID");
    return;
  }
  const button = shouldFollow ? el.followButton : el.unfollowButton;
  setBusy(button, true);
  try {
    await request(`/follows/${followeeId}`, { method: shouldFollow ? "POST" : "DELETE" });
    showToast(shouldFollow ? "已关注" : "已取关");
    setLastAction(shouldFollow ? "关注用户" : "取关用户");
    await loadFeed("latest", { quiet: true });
  } catch (error) {
    showToast(error.message);
    setLastAction(shouldFollow ? "关注失败" : "取关失败");
  } finally {
    setBusy(button, false);
  }
}

async function deletePost() {
  if (!state.menuContentId) return;
  const contentId = state.menuContentId;
  hidePostMenu();
  setBusy(el.deletePostButton, true);
  try {
    await request(`/posts/${contentId}`, { method: "DELETE" });
    state.items = state.items.filter((item) => String(item.content_id) !== String(contentId));
    renderFeed();
    showToast("已删除");
    setLastAction("删除动态");
  } catch (error) {
    showToast(error.message);
    setLastAction("删除失败");
  } finally {
    setBusy(el.deletePostButton, false);
  }
}

async function loadFeed(direction, options = {}) {
  if (!state.token) {
    showToast("请先登录");
    return;
  }
  if (state.feedLoading) return;

  const params = new URLSearchParams({ direction, limit: "20" });
  if (direction === "newer" && state.topCursor) params.set("cursor", state.topCursor);
  if (direction === "older" && state.bottomCursor) params.set("cursor", state.bottomCursor);

  const button = direction === "older" ? el.loadOlderButton : el.refreshButton;
  state.feedLoading = true;
  setBusy(button, true);
  el.refreshButton.classList.add("active");
  setHint(direction === "older" ? "正在加载更早内容" : "正在刷新", true);
  try {
    const data = await request(`/feed/timeline?${params.toString()}`, { method: "GET" });
    applyFeed(direction, data);
    renderFeed();
    updateCursorStatus();
    if (!options.quiet) {
      showToast(direction === "older" ? "已加载更早内容" : "已刷新");
    }
    setLastAction(direction === "older" ? "加载更早" : "刷新时间线");
    setHint("准备就绪");
  } catch (error) {
    showToast(error.message);
    setLastAction("读取失败");
    setHint("刷新失败");
  } finally {
    state.feedLoading = false;
    el.refreshButton.classList.remove("active");
    setBusy(button, false);
  }
}

function applyFeed(direction, data) {
  const incoming = data.items || [];
  if (direction === "older") {
    state.items = dedupe([...state.items, ...incoming]);
  } else if (direction === "newer") {
    state.items = dedupe([...incoming, ...state.items]);
  } else {
    state.items = incoming;
  }
  if (state.items.length) {
    state.topCursor = data.top_cursor || state.topCursor;
    state.bottomCursor = data.bottom_cursor || state.bottomCursor;
  } else {
    state.topCursor = "";
    state.bottomCursor = "";
  }
  state.hasMore = Boolean(data.has_more);
}

function dedupe(items) {
  const seen = new Set();
  return items.filter((item) => {
    const id = String(item.content_id);
    if (seen.has(id)) return false;
    seen.add(id);
    return true;
  });
}

function renderFeed() {
  hidePostMenu();
  if (!state.items.length) {
    el.feedList.innerHTML = `
      <div class="empty-view">
        <div>
          <strong>暂无动态</strong>
          <span>${state.token ? "关注的人发布后刷新即可看到" : "登录后刷新关注流"}</span>
        </div>
      </div>
    `;
    return;
  }

  el.feedList.innerHTML = state.items.map((item, index) => renderItem(item, index)).join("");
  bindMoreButtons();
  mountIcons();
}

function renderItem(item, index) {
  const content = item.content_text || "";
  const hasMedia = content.length > 18 && index % 3 === 1;
  const isMine = String(item.author_id) === String(state.userId);
  const authorLabel = `用户 ${item.author_id}`;
  return `
    <article class="feed-item ${hasMedia ? "has-media" : ""} ${isMine ? "is-mine" : ""}">
      <div class="feed-avatar" aria-hidden="true">${avatarArt(index)}</div>
      <div class="feed-main">
        <div class="feed-head">
          <h2 class="feed-author">${escapeHTML(authorLabel)}</h2>
          <button class="more-button" type="button" data-content-id="${escapeHTML(String(item.content_id))}" title="更多">
            <i data-lucide="more-horizontal"></i>
          </button>
        </div>
        <p class="post-body">${escapeHTML(content)}</p>
        <div class="post-media">
          <div class="media-frame">
            <div class="media-thumb"><i data-lucide="play"></i></div>
            <div class="media-title">${escapeHTML(content.slice(0, 36))}</div>
          </div>
        </div>
        <div class="feed-foot">
          <span>${relativeTime(item.publish_time)}</span>
          <span>·</span>
          <strong>关注流</strong>
          <span>·</span>
          <span>内容 ${escapeHTML(String(item.content_id))}</span>
        </div>
        <div class="like-row">
          <i data-lucide="heart"></i>
          <span>${escapeHTML(authorLabel)}</span>
        </div>
      </div>
    </article>
  `;
}

function bindMoreButtons() {
  document.querySelectorAll(".more-button").forEach((button) => {
    button.addEventListener("click", (event) => {
      const id = button.getAttribute("data-content-id") || "";
      showPostMenu(id, event.currentTarget);
    });
  });
}

function showPostMenu(contentId, anchor) {
  state.menuContentId = contentId;
  const rect = anchor.getBoundingClientRect();
  el.postMenu.style.left = `${Math.max(12, rect.right - 148)}px`;
  el.postMenu.style.top = `${rect.bottom + 8}px`;
  el.postMenu.classList.remove("hidden");
}

function hidePostMenu() {
  el.postMenu.classList.add("hidden");
}

function updateCursorStatus() {
  if (!state.topCursor && !state.bottomCursor) {
    el.cursorStatus.textContent = "Cursor -";
    return;
  }
  el.cursorStatus.textContent = state.hasMore ? "Cursor 可加载" : "Cursor 当前页";
}

function updateCounter() {
  el.postCounter.textContent = `${el.postText.value.length} / 2000`;
}

function handleWheel(event) {
  if (!state.token || state.feedLoading) return;
  const now = Date.now();
  if (now - state.wheelAt < 900) return;

  const nearBottom = el.stream.scrollTop + el.stream.clientHeight >= el.stream.scrollHeight - 8;
  if (event.deltaY < 0) {
    state.wheelAt = now;
    loadFeed("latest");
  } else if (event.deltaY > 0 && nearBottom) {
    state.wheelAt = now;
    loadFeed("older");
  }
}

function relativeTime(value) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return "-";
  const delta = Date.now() - date.getTime();
  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;
  if (delta < minute) return "刚刚";
  if (delta < hour) return `${Math.floor(delta / minute)}分钟前`;
  if (delta < day) return `${Math.floor(delta / hour)}小时前`;
  return new Intl.DateTimeFormat("zh-CN", {
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    hour12: false,
  }).format(date);
}

function avatarArt(index) {
  const colors = [
    ["#dbeafe", "#fef3c7", "#2d4c84"],
    ["#dcfce7", "#fee2e2", "#0f766e"],
    ["#fae8ff", "#e0f2fe", "#7c2d12"],
    ["#ede9fe", "#fef9c3", "#365314"],
  ][index % 4];
  return `
    <svg viewBox="0 0 74 74" role="img">
      <rect width="74" height="74" rx="10" fill="${colors[0]}"></rect>
      <circle cx="23" cy="22" r="12" fill="${colors[1]}"></circle>
      <circle cx="49" cy="20" r="10" fill="#fff"></circle>
      <path d="M12 58 C18 38, 34 36, 42 48 C48 57, 56 44, 66 54 L66 74 L12 74 Z" fill="${colors[2]}" opacity=".82"></path>
      <circle cx="31" cy="34" r="5" fill="#1f2429" opacity=".55"></circle>
    </svg>
  `;
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

el.loginTab.addEventListener("click", () => setMode("login"));
el.registerTab.addEventListener("click", () => setMode("register"));
el.authForm.addEventListener("submit", submitAuth);
el.logoutButton.addEventListener("click", () => {
  clearSession();
  showToast("已退出");
  setLastAction("退出登录");
});
el.postText.addEventListener("input", updateCounter);
el.publishButton.addEventListener("click", publishPost);
el.refreshButton.addEventListener("click", () => loadFeed("latest"));
el.cameraButton.addEventListener("click", () => el.postText.focus());
el.notifyButton.addEventListener("click", () => showToast(el.apiStatus.textContent));
el.loadOlderButton.addEventListener("click", () => loadFeed("older"));
el.followButton.addEventListener("click", () => followUser(true));
el.unfollowButton.addEventListener("click", () => followUser(false));
el.deletePostButton.addEventListener("click", deletePost);
el.stream.addEventListener("wheel", handleWheel, { passive: true });
document.addEventListener("click", (event) => {
  if (!el.postMenu.contains(event.target) && !event.target.closest(".more-button")) {
    hidePostMenu();
  }
});

setMode("login");
updateSessionView();
updateCounter();
renderFeed();
updateCursorStatus();
checkHealth();
mountIcons();
if (state.token) {
  loadFeed("latest", { quiet: true });
}
