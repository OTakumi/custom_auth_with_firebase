FROM node:20-alpine

# Alpine では標準リポジトリで Java 21 が提供されています
RUN apk add --no-cache openjdk21-jre-headless

# Firebase CLIのインストール
RUN npm install -g firebase-tools

WORKDIR /app

COPY firebase.json ./

# エミュレータの各ポート（Auth, Firestore, UI等）を公開
# 4000: UI, 8080: Firestore, 9099: Auth
EXPOSE 4000 8080 9099

# コンテナ起動時にエミュレータを開始
# --import/--export-on-exit をつけるとデータの永続化が可能
CMD ["firebase", "emulators:start", "--project", "demo-project", "--config", "firebase.json"]
