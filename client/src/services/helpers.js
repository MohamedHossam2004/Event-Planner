export const formatDate = (mongoDate) => {
  const date = new Date(mongoDate);
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
};

export const formatTime = (mongoDate) => {
  const date = new Date(mongoDate);
  return date.toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
  });
};