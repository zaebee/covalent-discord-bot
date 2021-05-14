import $http from '~/plugins/axios'
import { GetterTree, ActionTree, MutationTree } from 'vuex'

export const state = () => ({
  things: [] as string[],
  memes: [] as Object[],
  total: 0,
  pageSize: 9,
})

export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  total: state => state.total,
  hasNext: state => state.memes.length < state.total
}

export const mutations: MutationTree<RootState> = {
  SET_TOTAL: (state, newTotal: number) => (state.total = newTotal),
  INIT_MEMES: (state, memes: Object[]) => (state.memes = memes),
  MORE_MEMES: (state, memes: Object[]) => (state.memes = state.memes.concat(memes)),
}

export const actions: ActionTree<RootState, RootState> = {
  async fetchMemes({ commit }, payload) {
    let page =(payload && payload.page) || 0
    const memes = await $http.post('/memes/_search', {
      from: page * this.state.pageSize,
      size: this.state.pageSize,
      sort: {Timestamp: {order: 'desc'}}
    })
    console.log(memes.data.hits.hits)
    if (page == 0) {
      commit('INIT_MEMES', memes.data.hits.hits)
    } else {
      commit('MORE_MEMES', memes.data.hits.hits)
    }
    commit('SET_TOTAL', memes.data.hits.total)
  },
  async fetchTopMemes({ commit }, payload) {
    let page =(payload && payload.page) || 0
    const memes = await $http.post('/memes/_search', {
      query: {
        bool: {
          filter: [{exists: {field: 'Reactions'}}]
        },
      },
      from: page * this.state.pageSize,
      size: this.state.pageSize,
      sort: {'Reactions.keyword': {order: 'desc'}}
    })
    console.log(memes.data.hits.hits)
    if (page == 0) {
      commit('INIT_MEMES', memes.data.hits.hits)
    } else {
      commit('MORE_MEMES', memes.data.hits.hits)
    }
    commit('SET_TOTAL', memes.data.hits.total)
  },
}
