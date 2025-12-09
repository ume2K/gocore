import { initStickyHeader } from './components/header.js';

document.addEventListener('DOMContentLoaded', () => {
    initStickyHeader();
});

document.addEventListener("DOMContentLoaded", function () {
    const currentPath = window.location.pathname;

    /**
     * Hilfsfunktion um Links zu aktivieren
     * @param {string} selector - Welcher Link soll gesucht werden?
     * @param {string} activeClass - Welche Klasse soll hinzugefügt werden?
     */
    const highlightActiveLink = (selector, activeClass) => {
        const links = document.querySelectorAll(selector);

        links.forEach((link) => {
            // 1. Reset: Alte Klassen entfernen (wichtig, da im HTML "is-active" hardcoded ist)
            link.classList.remove(activeClass);
            link.setAttribute("aria-selected", "false");

            // 2. Check: Stimmt der Pfad überein?
            // Wir nutzen getAttribute("href"), damit wir exakt das bekommen, was im HTML steht
            if (link.getAttribute("href") === currentPath) {
                link.classList.add(activeClass);
                link.setAttribute("aria-selected", "true");
            }
        });
    };

    // --- Desktop Navigation ---
    // Hier fügen wir die Klasse "active" hinzu (für den Unterstrich-Effekt)
    highlightActiveLink(".navigations .link", "active");

    // --- Mobile App Navigation ---
    // Hier fügen wir die Klasse "is-active" hinzu (für das Einfärben & Icon-Effekt)
    highlightActiveLink(".app-nav-btn", "active");
});