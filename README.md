# Deep Art

**Deep Art** is a full-stack e-commerce-like artwork platform, allowing users to upload their artwork and set a price for others to favorite, buy, and download.

## Features
- **User:** View all available arts. Favorite, buy, and download artworks.
- **Creator:** Manage owned artworks with full CRUD functionality.

## Technologies Used
- **Frontend:** HTMX, Alpine.js, TailwindCSS
- **Backend:** Golang
- **Deployment:** Docker, Google Cloud Run

## Setup Instructions

### Prerequisites
- [Golang](https://go.dev/dl/)
- [Air](https://github.com/air-verse/air)
- [Templ](https://templ.guide/quick-start/installation)
- [Node](https://nodejs.org/en/download)
- Set up create Github OAuth.
- Set up a [Google Cloud Project](https://cloud.google.com/docs/project), create Bucket Storage and OAuth.

### Local Development

1. Clone the repository.
   ```bash
   git clone https://github.com/DeepAung/deep-art.git
   cd deep-art
   ```
2. Create `.env.dev` with template from `.env.example`
3. Run migration
   ```bash
   make migrate.up
   ```
4. Run these 3 commands separately.
   ```bash
   make air
   ```
   ```bash
   make tailwind
   ```
   ```bash
   make templ
   ```
5. Access the application at [http://localhost:3000](http://localhost:3000).

## Live Demo

Experience the live application at [https://deep-art-prod-796109602795.asia-southeast1.run.app](https://deep-art-prod-796109602795.asia-southeast1.run.app).
