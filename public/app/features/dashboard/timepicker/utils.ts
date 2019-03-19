const unitRegex = /(\d*\.?\d+)\s?(.*)/;

export const DEFAULT_MIN_INTERVAL_SECONDS = 60;

export const isAboveMinThresholdSeconds = (interval: string, threshold = DEFAULT_MIN_INTERVAL_SECONDS) => {
  // If the interval is in seconds, check to see if it's less than 60s
  const parts = unitRegex.exec(interval);

  if (parts[2] === 's' || parts[2] === 'seconds') {
    const parts = interval.split('s');
    const seconds = parseInt(parts[0], 10);

    if (isNaN(seconds)) {
      return false;
    }

    return seconds >= threshold;
  }

  return true;
};

export const createRefreshIntervalOptions = (intervals: string[]) => {
  return intervals.filter((interval: string) => isAboveMinThresholdSeconds(interval)).map((interval: string) => {
    return { text: interval, value: interval };
  });
};
