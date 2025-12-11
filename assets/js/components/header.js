export function initStickyHeader() {
  const header = document.getElementById("gmuende-header");
  const arrow = document.getElementById("arrow");
  if (!header) return;

  // Konfiguration
  const SCROLL_THRESHOLD = 50;
  const MOUSE_THRESHOLD = 100;

  // State
  let isScrolled = window.scrollY > SCROLL_THRESHOLD;
  let isMouseTop = false;
  let rafPending = false;

  const render = () => {
    const shouldHighlight = isMouseTop;
    const shouldCollapse = isScrolled && !isMouseTop;

    header.classList.toggle("highlighted", shouldHighlight);
    header.classList.toggle("collapsed", shouldCollapse);
    arrow?.classList.toggle("hide", isScrolled);    

    rafPending = false;
  };

  const requestUpdate = () => {
    if (!rafPending) {
      requestAnimationFrame(render);
      rafPending = true;
    }
  };

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

  window.addEventListener(
    "mousemove",
    (event) => {
      const currentlyTop = event.clientY <= MOUSE_THRESHOLD;
      if (isMouseTop !== currentlyTop) {
        isMouseTop = currentlyTop;
        requestUpdate();
      }
    },
    { passive: true }
  );
  
  requestUpdate();
}