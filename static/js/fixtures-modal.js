(function () {
    console.log("HERE");
    /** @type HTMLDivElement */
    const modalWrapper = document.getElementById("modal-wrapper");
    /** @type HTMLDivElement */
    const modalOverlay = document.getElementById("modal-overlay");

    function hideModalContent() {
        modalOverlay.classList.remove("opacity-20");
        modalOverlay.classList.add("opacity-0");
        const modalContent = document.getElementById("modal-content");
        modalContent.classList.remove("opacity-100");
        modalContent.classList.add("opacity-0");
        setTimeout(() => {
            modalWrapper.classList.remove("flex");
            modalWrapper.classList.add("hidden");
            modalContent.innerHTML = "";
        }, 150);
    }

    /**
     * @param {KeyboardEvent} event
     */
    function onKeyDownModalOverlay(event) {
        if (event.key !== "Escape") return;
        document.removeEventListener("keydown", onKeyDownModalOverlay);
        hideModalContent();
    }

    modalWrapper.addEventListener("htmx:beforeSwap", function () {
        this.classList.remove("hidden");
        this.classList.add("flex");
        const modalContent = document.getElementById("modal-content");
        requestAnimationFrame(() => {
            modalOverlay.classList.remove("opacity-0");
            modalOverlay.classList.add("opacity-20");
            modalContent.classList.remove("opacity-0");
            modalContent.classList.add("opacity-100");
        });
        modalOverlay.addEventListener("click", hideModalContent, {
            once: true,
        });
        document.addEventListener("keydown", onKeyDownModalOverlay);
    });
})();
