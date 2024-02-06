function init() {
  Alpine.store("data", { passport: {} });
}

async function login(e) {
  try {
    const res = await fetch("/api/v1/users/login", {
      method: "POST",
      body: new FormData(e.target),
    });
    Alpine.store("data").passport = await res.json();
  } catch (e) {
    console.log("error: ", e);
  }
}

async function connectGoogle() {
  try {
    const accessToken = Alpine.store("data").passport.token.access_token;
    const res = await fetch("/api/v1/users/google/connect?provider=google", {
      method: "GET",
      headers: {
        "Authorization": "Bearer " + accessToken,
      },
    });

    console.log(res);
  } catch (e) {
    console.log("error: ", e);
  }
}
