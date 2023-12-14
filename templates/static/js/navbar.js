const navbar = document.getElementById("navbar");
window.onscroll = (e) => {
    if (window.scrollY > 100) {
        navbar.classList.replace("nav-transparent", "nav-colored");
    } else {
        navbar.classList.replace("nav-colored", "nav-transparent");
    }
};