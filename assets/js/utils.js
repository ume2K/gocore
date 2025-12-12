export function initScrollReveal(selector = '.reveal', activeClass = 'active') {
    const isMobile = window.matchMedia("(max-width: 900px)").matches;

    const bottomMargin = isMobile ? "-120px" : "-50px";

    const observerOptions = {
        threshold: 0.15,
        rootMargin: `0px 0px ${bottomMargin} 0px`
    };

    const observer = new IntersectionObserver((entries, obs) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add(activeClass);
                obs.unobserve(entry.target);
            }
        });
    }, observerOptions);

    const elements = document.querySelectorAll(selector);
    
    if (elements.length > 0) {
        elements.forEach(el => observer.observe(el));
    }
}