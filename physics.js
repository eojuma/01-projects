function getAcceleration(accel) {
  if (accel.f && accel.m) {
    return force / mass;
  } else if (accel.Δv && accel.Δt) {
    return Δv / Δt;
  } else if (accel.d && accel.t) {
    return (2 * d) / (t * t);
  } else {
    return "impossible";
  }
}
