// Routing für die Bottom-App-Navigation mit History API
// Erwartet:
// - Buttons: .app-nav-btn[data-route]  (z.B. data-route="/besuch")
// - Sections: [data-tab="<key>"]      (z.B. data-tab="besuch")

const routes = {
  "/": "home",
  "/home": "home",
  "/besuch": "besuch",
  "/tiere": "tiere",
  "/eventraum": "eventraum"
};

const defaultRoute = "/";

function normalizePath(pathname) {
  const path = pathname.replace(/\/+$/, "") || "/";
  return routes[path] ? path : defaultRoute;
}

function setActiveNav(route) {
  const selector = ".app-nav-btn, .navigations a";
  
  document.querySelectorAll(selector).forEach((el) => {
    const targetPath = el.dataset.route;
    const isActive = targetPath === route;
    
    el.classList.toggle("active", isActive);
    
    if (el.hasAttribute("aria-selected")) {
      el.setAttribute("aria-selected", isActive ? "true" : "false");
    }
  });
}

function setActiveSection(tabKey) {
  document.querySelectorAll("[data-tab]").forEach((section) => {
    const isActive = section.dataset.tab === tabKey;
    section.classList.toggle("active", isActive);
  });
}

function navigateTo(route, push = true) {
  const tabKey = routes[route];
  if (!tabKey) return;

  if (push) {
    window.history.pushState({ route }, "", route);
  }

  setActiveNav(route);
  setActiveSection(tabKey);
}

function handlePopState(event) {
  const route = event.state?.route || normalizePath(window.location.pathname);
  navigateTo(route, false);
}

function handleNavClick(event) {
  const btn = event.currentTarget;
  const route = btn.dataset.route;
  if (!route) return;

  event.preventDefault();
  navigateTo(route, true);
}

export function initRouting() {
  // Click-Handler für Nav-Buttons
  document.querySelectorAll(".app-nav-btn").forEach((btn) => {
    btn.addEventListener("click", handleNavClick);
  });

  // Initialen Zustand anhand der URL setzen
  const initialRoute = normalizePath(window.location.pathname);
  navigateTo(initialRoute, false);

  // Browser-Back/Forward
  window.addEventListener("popstate", handlePopState);
}