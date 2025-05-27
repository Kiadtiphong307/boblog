FROM node:20-alpine

WORKDIR /app

# Install dependencies first (better caching)
COPY package*.json ./
RUN npm install

# Copy the rest of the application
COPY . .

# Development mode - no build step
EXPOSE 3000

# Start the development server with hot reload
CMD ["npm", "run", "dev"]
