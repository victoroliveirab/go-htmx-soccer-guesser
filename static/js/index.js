/** @type HTMLDivElement */
const mobileMenuButton = document.getElementById("mobile-menu-button");
/** @type HTMLDivElement */
const mobileMenu = document.getElementById("mobile-menu");

mobileMenuButton.addEventListener("click", () => {
    const svgs = mobileMenuButton.getElementsByTagName("svg");
    for (const svg of svgs) {
        svg.classList.toggle("hidden");
    }
    mobileMenu.children[0].classList.toggle("hidden");
});
