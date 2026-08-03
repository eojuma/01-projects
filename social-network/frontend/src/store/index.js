import { fetchGroups } from "../services/groups.js"
import { fetchGroupPosts } from "../services/groups.js"
import { fetchGroupComments } from "../services/groups.js"

const state = {
  groups: [],
  groupPosts: [],  
  groupComments: [],
}

const mutations = {
  SET_GROUPS(state, groups) {
    state.groups = groups
  },
  SET_GROUP_POSTS(state, posts) {
    state.groupPosts = posts
  },
  SET_GROUP_COMMENTS(state, comments) {
    state.groupComments = comments
  },
}

const actions = {
  async loadGroups({ commit }, { page, limit, search } = {}) {
    const groups = await fetchGroups(page, limit, search)
    commit("SET_GROUPS", groups)
  },
  async loadGroupPosts({ commit }, { groupID, page, limit } = {}) {
    const posts = await fetchGroupPosts(groupID, page, limit)
    commit("SET_GROUP_POSTS", posts)
  },
   async loadGroupComments({ commit }, { groupID, postID } = {}) {
    const comments = await fetchGroupComments(groupID, postID)
    commit("SET_GROUP_COMMENTS", comments)
  },
}

export default {
  state,
  mutations,
  actions,
}

