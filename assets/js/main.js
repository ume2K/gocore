import { initRouting } from './components/nav.js';
import { initStickyHeader } from './components/header.js';

document.addEventListener('DOMContentLoaded', () => {
    initRouting();
    initStickyHeader();
});