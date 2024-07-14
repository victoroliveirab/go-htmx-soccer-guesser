/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["components/**/*.html", "views/**/*.html", "static/js/**/*.js"],
    theme: {
        extend: {
            spacing: {
                128: "32rem",
            },
        },
    },
};
