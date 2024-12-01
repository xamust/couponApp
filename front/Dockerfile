FROM node:18
WORKDIR /app/admin
RUN ls
COPY package*.json ./
RUN npm install
COPY . .
EXPOSE 3000
WORKDIR /app/admin
CMD ["npm", "start"]
