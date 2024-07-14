(function () {
    const checkboxes = document.querySelectorAll(
        "#ranking-content input[type=checkbox]"
    );
    checkboxes.forEach((checkbox, index) => {
        const divId = `ranking-content-wrapper-${index}`;
        const div = document.getElementById(divId);
        checkbox.addEventListener("change", (e) => {
            const isChecked = e.target.checked;
            if (isChecked) {
                div.classList.remove("max-h-0");
                div.classList.add("max-h-64");
                setTimeout(() => {
                    div.classList.remove("overflow-y-hidden");
                    div.classList.add("overflow-y-auto");
                }, 500);
            } else {
                div.classList.add("max-h-0");
                div.classList.add("overflow-y-hidden");
                div.classList.remove("overflow-y-auto");
                div.classList.remove("max-h-64");
            }
        });
    });
})();

function rankingTableIsCheckboxOn() {
    return this.checked;
}
