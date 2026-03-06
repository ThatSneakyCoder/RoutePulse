FROM node:20-alpine

WORKDIR /app

# install deps
COPY package*.json ./
RUN npm ci

# copy source
COPY . .

EXPOSE 5173

# vite dev server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]