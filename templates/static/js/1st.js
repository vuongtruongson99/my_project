const darkSwitch = document.querySelector(".dark-switch");
const lightSwitch = document.querySelector(".light-switch");
const body = document.querySelector("body");
const section1 = document.querySelector(".section-1");
const btn_hero = document.querySelector(".main-btn");

if (localStorage.getItem("darkMode")) {
    body.classList.add("dark-mode");
    section1.classList.add("dark-bg");
    darkSwitch.classList.remove("active");
    lightSwitch.classList.add("active");
    btn_hero.classList.add("dark");
    localStorage.setItem("darkMode", "true");
}

darkSwitch.addEventListener("click", () => {
    body.classList.add("dark-mode");
    section1.classList.add("dark-bg");
    darkSwitch.classList.remove("active");
    lightSwitch.classList.add("active");
    btn_hero.classList.add("dark");
    localStorage.setItem("darkMode", "true");
})

lightSwitch.addEventListener("click", () => {
    body.classList.remove("dark-mode");
    section1.classList.remove("dark-bg");
    darkSwitch.classList.add("active");
    lightSwitch.classList.remove("active");
    btn_hero.classList.remove("dark");
    localStorage.removeItem("darkMode");
})