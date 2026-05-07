function getAcceleration(accel) {
  if (accel.f && accel.m) {
    return f / m;
  } else if (accel.Δv && accel.Δt) {
    return Δv / Δt;
  } else if (accel.d && accel.t) {
    return (2 * d) / (t * t);
  } else {
    return "impossible";
  }
}
