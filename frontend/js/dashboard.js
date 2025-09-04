// js/dashboard.js
const API_BASE_URL = 'https://appointy-task-service.onrender.com'; // Replace with your backend service URL when deploying
const token = localStorage.getItem('jwt_token');

// Elements
const logoutBtn = document.getElementById('logout-btn');
const createLinkForm = document.getElementById('create-link-form');
const longUrlInput = document.getElementById('long-url');
const linksContainer = document.getElementById('links-container');

// --- Main execution on page load ---
document.addEventListener('DOMContentLoaded', () => {
    // This is a "route guard". If no token, redirect to login.
    if (!token) {
        window.location.href = 'login.html';
        return;
    }
    fetchLinks();
});

// --- Event Listeners ---
logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('jwt_token');
    window.location.href = 'login.html';
});

createLinkForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const longUrl = longUrlInput.value;
    try {
        // --- THIS LINE IS NOW CORRECTED ---
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
    links.forEach(link => {
        const linkCard = document.createElement('div');
        linkCard.className = 'card link-card';
        // Use the backend URL for the short link to ensure it works when deployed
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
        // Fetch analytics for this specific link
        fetchAnalytics(link.short_id);
    });
}

async function fetchAnalytics(linkId) {
    try {
        const response = await fetch(`${API_BASE_URL}/analytics/${linkId}`, {
            method: 'GET',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (!response.ok) return; // Fail silently if analytics aren't available
        const data = await response.json();
        const clicksEl = document.getElementById(`clicks-${linkId}`);
        if(clicksEl) {
            clicksEl.textContent = data.total_clicks;
        }
    } catch (error) {
        console.error(`Could not fetch analytics for ${linkId}`, error);
    }
}