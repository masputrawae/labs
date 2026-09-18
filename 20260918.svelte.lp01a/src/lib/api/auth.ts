export const isAuth = async (): Promise<boolean> => {
  try {
    const res = fetch("/api/me", {
      method: "GET",
      credentials: "include",
    });
    if ((await res).ok) {
      return true;
    } else {
      return false;
    }
  } catch (error) {
    console.log(error);
    return false;
  }
};
