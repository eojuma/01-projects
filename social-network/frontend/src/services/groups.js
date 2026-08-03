const API_URL = "";

export async function createGroup(title, description) {
  const res = await fetch(`${API_URL}/api/groups`, {
    method: "POST",
    credentials: "include", // IMPORTANT for session cookie
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      title,
      description,
    }),
  });

  if (!res.ok) {
    const err = await res.text();
    throw new Error(err || "Failed to create group");
  }

  return res.json();
}


export async function fetchGroups(page = 1, limit = 20, search = "") {
  const params = new URLSearchParams({ page, limit, search })

  const res = await fetch(`${API_URL}/api/groups?${params}`, {
    method: "GET",
    credentials: "include",
  })

  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to fetch groups")
  }

  return res.json()
}

export async function inviteUser(groupID, userID) {
  const res = await fetch(`${API_URL}/api/groups/${groupID}/invite`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ user_id: userID }),
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to invite user")
  }
  return res.json()
}

export async function acceptInvite(groupID, inviteID) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/invites/${inviteID}/accept`,
    {
      method: "POST",
      credentials: "include",
    }
  )
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to accept invite")
  }
  return res.json()
}

export async function declineInvite(groupID, inviteID) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/invites/${inviteID}/decline`,
    {
      method: "POST",
      credentials: "include",
    }
  )
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to decline invite")
  }
  return res.json()
}

export async function requestToJoin(groupID) {
  const res = await fetch(`${API_URL}/api/groups/${groupID}/request`, {
    method: "POST",
    credentials: "include",
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to send join request")
  }
  return res.json()
}

export async function approveRequest(groupID, requestID) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/requests/${requestID}/approve`,
    {
      method: "POST",
      credentials: "include",
    }
  )
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to approve request")
  }
  return res.json()
}

export async function declineRequest(groupID, requestID) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/requests/${requestID}/decline`,
    {
      method: "POST",
      credentials: "include",
    }
  )
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to decline request")
  }
  return res.json()
}

export async function createGroupPost(groupID, content) {
  const res = await fetch(`${API_URL}/api/groups/${groupID}/posts`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ content }),
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to create group post")
  }
  return res.json()
}

export async function fetchGroupPosts(groupID, page = 1, limit = 20) {
  const params = new URLSearchParams({ page, limit })

  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/posts?${params}`,
    {
      method: "GET",
      credentials: "include",
    }
  )

  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to fetch group posts")
  }

  return res.json()
}

export async function createGroupComment(groupID, postID, content) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/posts/${postID}/comments`,
    {
      method: "POST",
      credentials: "include",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ content }),
    }
  )

  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to create comment")
  }

  return res.json()
}

export async function fetchGroupComments(groupID, postID) {
  const res = await fetch(
    `${API_URL}/api/groups/${groupID}/posts/${postID}/comments`,
    {
      method: "GET",
      credentials: "include",
    }
  )

  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to fetch comments")
  }

  return res.json()
}

export async function createGroupEvent(groupID,title,description,eventDate){
  const res =await fetch(`${API_URL}/api/groups/${groupID}/events`,{
    method:"POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
    title,
    description,
    event_date: eventDate,
    }),
  })

  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to create event")
  }

  return res.json()
}

export async function fetchGroupEvents(groupID) {
  const res = await fetch(`${API_URL}/api/groups/${groupID}/events`, {
    method: "GET",
    credentials: "include",
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to fetch events")
  }
  return res.json()
}

export async function respondToEvent(eventID, response) {
  const res = await fetch(`${API_URL}/api/events/${eventID}/respond`, {
    method: "POST",
    credentials: "include",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ response }),
  })
  if (!res.ok) {
    const err = await res.text()
    throw new Error(err || "Failed to respond to event")
  }
  return res.json()
}