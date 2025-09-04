// js/dashboard.js
const API_BASE_URL = 'https://appointy-task-service.onrender.com'; // Your backend service URL
const token = localStorage.getItem('jwt_token');

// Elements
const logoutBtn = document.getElementById('logout-btn');
const createLinkForm = document.getElementById('create-link-form');
const longUrlInput = document.getElementById('long-url');
const linksContainer = document.getElementById('links-container');

// A variable to hold our polling interval so we can stop it if needed.
let analyticsInterval;

// --- Main execution on page load ---
document.addEventListener('DOMContentLoaded', () => {
    if (!token) {
        window.location.href = 'index.html'; // Changed from login.html
        return;
    }
    fetchLinks();
});

// --- Event Listeners ---
logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('jwt_token');
    clearInterval(analyticsInterval); // Stop polling when logging out
    window.location.href = 'index.html'; // Changed from login.html
});

createLinkForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const longUrl = longUrlInput.value;
    try {
        const response = await fetch(`${API_BASE_URL}/shorten`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ url: longUrl })
        });
        if (!response.ok) throw new Error('Failed to create link.');
        longUrlInput.value = '';
        fetchLinks(); // Refresh the list of links
    } catch (error) {
        alert(error.message);
    }
});

// --- Functions ---
async function fetchLinks() {
    // Stop any previous polling before we re-render the list
    if (analyticsInterval) clearInterval(analyticsInterval);

    try {
        const response = await fetch(`${API_BASE_URL}/links`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) throw new Error('Could not fetch links.');
        const links = await response.json();
        renderLinks(links || []);
    } catch (error) {
        console.error(error);
    }
}

function renderLinks(links) {
    linksContainer.innerHTML = '';
    if (links.length === 0) {
        linksContainer.innerHTML = '<p>You haven\'t created any links yet.</p>';
        return;
    }

    const linkIds = []; // Keep track of all link IDs on the page
    links.forEach(link => {
        linkIds.push(link.short_id);
        const linkCard = document.createElement('div');
        linkCard.className = 'card link-card';
        const shortUrl = `${API_BASE_URL}/r/${link.short_id}`;
        linkCard.innerHTML = `
            <div class="link-info">
                <a href="${shortUrl}" target="_blank" class="short-url">${shortUrl}</a>
                <p class="original-url">${link.original_url}</p>
            </div>
            <div class="link-analytics">
                <h3 id="clicks-${link.short_id}">-</h3>
                <p>Total Clicks</p>
            </div>
        `;
        linksContainer.appendChild(linkCard);
        fetchAnalytics(link.short_id); // Fetch initial analytics
    });

    // --- NEW POLLING LOGIC ---
    // After rendering all links, start an interval to refresh their analytics every 5 seconds.
    analyticsInterval = setInterval(() => {
        console.log("Polling for new analytics...");
        linkIds.forEach(id => fetchAnalytics(id));
    }, 2000);
}

async function fetchAnalytics(linkId) {
    try {
        const response = await fetch(`${API_BASE_URL}/analytics/${linkId}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) return;
        const data = await response.json();
        const clicksEl = document.getElementById(`clicks-${linkId}`);
        if(clicksEl) {
            clicksEl.textContent = data.total_clicks;
        }
    } catch (error) {
        // Don't alert for polling errors, just log them.
        console.error(`Could not fetch analytics for ${linkId}`, error);
    }
}