import { initScrollReveal } from './utils.js';
import { initStickyHeader } from './components/header.js';

document.addEventListener('DOMContentLoaded', () => {
    initScrollReveal();
    initStickyHeader();
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
    
    if (scrollPosition > hero.offsetHeight) return;

    const speed = 0.5;
    const offset = scrollPosition * speed;

    hero.style.setProperty('--parallax-offset', `${offset}px`);
});