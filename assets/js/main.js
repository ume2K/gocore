import { initScrollReveal } from './utils.js';
import { initDesktopNav } from './components/nav.js';

document.addEventListener('DOMContentLoaded', () => {
    initScrollReveal();
    initDesktopNav();
});

document.addEventListener("DOMContentLoaded", function () {
    const currentPath = window.location.pathname;

    const highlightActiveLink = (selector, activeClass) => {
        const links = document.querySelectorAll(selector);

        links.forEach((link) => {
            link.classList.remove(activeClass);
            link.setAttribute("aria-selected", "false");

            if (link.getAttribute("href") === currentPath) {
                link.classList.add(activeClass);
                link.setAttribute("aria-selected", "true");
            }
        });
    };

    highlightActiveLink(".navigations .link", "active");
    highlightActiveLink(".app-nav-btn", "active");
});

document.addEventListener('scroll', () => {
    const hero = document.querySelector('.hero');
    if (!hero) return;

    const scrollPosition = window.scrollY;
    const heroHeight = hero.offsetHeight;

    if (scrollPosition > (heroHeight + 100)) return;

    const speed = 0.5;
    const offset = scrollPosition * speed; 
    hero.style.setProperty('--parallax-offset', `${offset}px`);

    // Konfiguration
    const blurFactor = 60; // Höherer Wert = langsamerer Blur-Anstieg
    const maxBlur = 10;    // Maximaler Blur in Pixeln (z.B. 15px)

    const blurAmount = Math.min(scrollPosition / blurFactor, maxBlur);

    hero.style.setProperty('--parallax-blur', `${blurAmount}px`);
});