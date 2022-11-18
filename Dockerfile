FROM ubuntu
EXPOSE 3000
RUN ["apt-get", "update"]
RUN ["apt-get", "-y", "install", "youtube-dl"]
COPY youtube_subs.html .
COPY youtubeSubtitles .
CMD ["./youtubeSubtitles"]
