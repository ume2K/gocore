export function initStickyHeader() {
  const header = document.getElementById("gmuende-header");
  const arrow = document.getElementById("arrow");
  if (!header && !arrow) return;

  // Konfiguration
  const SCROLL_THRESHOLD = 50;
  const MOUSE_THRESHOLD = 100;

  // State
  let isScrolled = window.scrollY > SCROLL_THRESHOLD;
  let isMouseTop = false;
  let rafPending = false;

  // Die eigentliche DOM-Manipulation (Render)
  const render = () => {
    // Logik:
    // 1. Highlighted: Immer wenn die Maus oben ist.
    // 2. Collapsed: Nur wenn gescrollt UND die Maus NICHT oben ist.
    const shouldHighlight = isMouseTop;
    const shouldCollapse = isScrolled && !isMouseTop;

    header.classList.toggle("highlighted", shouldHighlight);
    header.classList.toggle("collapsed", shouldCollapse);
    arrow.classList.toggle("hide", isScrolled);    

    rafPending = false;
  };

  // Hilfsfunktion: Fordert ein Update beim nächsten Frame an
  const requestUpdate = () => {
    if (!rafPending) {
      requestAnimationFrame(render);
      rafPending = true;
    }
  };

  // Event Listener: Scroll
  window.addEventListener(
    "scroll",
    () => {
      const currentlyScrolled = window.scrollY > SCROLL_THRESHOLD;
      if (isScrolled !== currentlyScrolled) {
        isScrolled = currentlyScrolled;
        requestUpdate();
      }
    },
    { passive: true }
  );

  // Event Listener: Mousemove
  window.addEventListener(
    "mousemove",
    (event) => {
      const currentlyTop = event.clientY <= MOUSE_THRESHOLD;
      // Nur updaten, wenn sich der Status wirklich geändert hat (Performance!)
      if (isMouseTop !== currentlyTop) {
        isMouseTop = currentlyTop;
        requestUpdate();
      }
    },
    { passive: true }
  );
  
  // Initialer Aufruf, um Status beim Laden zu setzen
  requestUpdate();
}