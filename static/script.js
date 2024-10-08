window.addEventListener('scroll', function() {
    var header = document.querySelector('header');
    var scrollPosition = window.scrollY;

    if (scrollPosition > 50) {
        header.classList.add('shrink');
    } else {
        header.classList.remove('shrink');
    }
});

// Count Animation
function animateCount(id, start, end, duration) {
    let obj = document.getElementById(id),
        current = start,
        range = end - start,
        increment = end > start ? 1 : -1,
        stepTime = Math.abs(Math.floor(duration / range)),
        timer = setInterval(function() {
            current += increment;
            obj.textContent = current;
            if (current == end) {
                clearInterval(timer);
            }
        }, stepTime);
}

document.addEventListener("DOMContentLoaded", function() {
    // Trigger animations on scroll
    const sections = document.querySelectorAll('.section');
    const heroSection = document.querySelector('.hero');
    
    const options = {
        threshold: 0.5
    };
    
    const observer = new IntersectionObserver(function(entries, observer) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                if (entry.target === heroSection) {
                    entry.target.classList.add('hero-fade');
                } else {
                    entry.target.classList.add('animate-pop');
                }
                observer.unobserve(entry.target);
            }
        });
    }, options);

    sections.forEach(section => {
        observer.observe(section);
    });

    // Start the counting animation when counts section is in view
    const countsSection = document.querySelector('#counts');
    const countsObserver = new IntersectionObserver(function(entries) {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                animateCount('clients-count', 0, 300, 2000);
                animateCount('projects-count', 0, 550, 2500);
                animateCount('employees-count', 0, 400, 2000);
                animateCount('products-count', 0, 500, 3000);
                countsObserver.unobserve(countsSection);
            }
        });
    }, { threshold: 0.5 });

    countsObserver.observe(countsSection);
});
