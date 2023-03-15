FROM node:19.4.0
COPY package-lock.json .
COPY package.json .
RUN npm install
COPY . .
CMD npm start