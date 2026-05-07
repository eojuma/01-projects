function getAcceleration(accel) {
  if (accel.f && accel.m) {
    return accel.f / accel.m;
  } else if (accel.Δv && accel.Δt) {
    return accel.Δv / accel.Δt;
  } else if (accel.d && accel.t) {
    return (2 * accel.d) / (accel.t * accel.t);
  } else {
    return "impossible";
  }
}
