const API_BASE_URL = 'http://localhost:3000';
const token = localStorage.getItem('jwt_token');

const logoutBtn = document.getElementById('logout-btn');
const createLinkForm = document.getElementById('create-link-form');
const longUrlInput = document.getElementById('long-url');
const linksContainer = document.getElementById('links-container');

document.addEventListener('DOMContentLoaded', () => {
    if (!token) {
        window.location.href = 'index.html';
        return;
    }
    fetchLinks();
    connectWebSocket();
});

logoutBtn.addEventListener('click', () => {
    localStorage.removeItem('jwt_token');
    window.location.href = 'index.html';
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
        fetchLinks();
    } catch (error) {
        alert(error.message);
    }
});

function connectWebSocket() {
    const wsUrl = API_BASE_URL.replace(/^http/, 'ws');
    const socket = new WebSocket(`${wsUrl}/ws?token=${token}`);

    socket.onopen = () => {
        console.log('WebSocket connection established.');
    };

    socket.onmessage = (event) => {
        console.log('Received message from server:', event.data);
        const message = JSON.parse(event.data);
        
        const clicksEl = document.getElementById(`clicks-${message.link_id}`);
        if (clicksEl) {
            clicksEl.textContent = message.click_count;
        }
    };

    socket.onclose = () => {
        console.log('WebSocket connection closed. Attempting to reconnect in 5 seconds...');
        setTimeout(connectWebSocket, 5000);
    };

    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

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
        fetchAnalytics(link.short_id);
    });
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
        console.error(`Could not fetch analytics for ${linkId}`, error);
    }
}