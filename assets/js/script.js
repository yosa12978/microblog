function toDateString_(date) {
  const dt = new Date(date);
  const options = {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    timeZone: Intl.DateTimeFormat().resolvedOptions().timeZone,
    hour12: false,
  };

  const dateString = dt.toLocaleString("default", options);
  const parts = dateString.split(" ");
  const timeParts = parts[parts.length - 1].split(":");
  return `${parts[0]} ${parts[1]} ${parts[2]} ${timeParts[0]}:${timeParts[1]}`;
}

function renderMarkdown(md) {
  return marked.parse(md);
}
