FROM node:20-alpine

WORKDIR /app

# install deps
COPY frontend/web/package*.json ./
RUN npm install

# copy source
COPY frontend/web ./

EXPOSE 5173

# vite dev server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]