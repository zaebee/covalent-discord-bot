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
  hasNext: state => state.memes.length < state.total,
}

export const mutations: MutationTree<RootState> = {
  SET_TOTAL: (state, newTotal: number) => (state.total = newTotal),
  INIT_MEMES: (state, memes: Object[]) => (state.memes = memes),
  MORE_MEMES: (state, memes: Object[]) =>
    (state.memes = state.memes.concat(memes)),
}

export const actions: ActionTree<RootState, RootState> = {
  async fetchMemes({ commit }, payload) {
    const { page = 0, userName } = payload
    const params: any = {
      from: page * this.state.pageSize,
      size: this.state.pageSize,
      sort: { Timestamp: { order: 'desc' } },
    }
    if (userName) {
      params.query = {
        match: {
          'Author.Username': userName,
        },
      }
    }

    const memes = await $http.post('/memes/_search', params)
    console.log(memes.data.hits.hits)
    if (page === 0) {
      commit('INIT_MEMES', memes.data.hits.hits)
    } else {
      commit('MORE_MEMES', memes.data.hits.hits)
    }
    commit('SET_TOTAL', memes.data.hits.total)
  },
  async fetchTopMemes({ commit }, payload) {
    const { page = 0, userName } = payload
    const params: any = {
      query: {
        bool: {
          filter: [{ exists: { field: 'Reactions' } }],
        },
      },
      from: page * this.state.pageSize,
      size: this.state.pageSize,
      sort: { 'Rating': { order: 'desc' } },
    }

    if (userName) {
      Object.assign(params.query.bool, {
        must: [
          {
            match: {
              'Author.Username': userName,
            },
          },
        ],
      })
    }

    const memes = await $http.post('/memes/_search', params)
    console.log(memes.data.hits.hits)
    if (page === 0) {
      commit('INIT_MEMES', memes.data.hits.hits)
    } else {
      commit('MORE_MEMES', memes.data.hits.hits)
    }
    commit('SET_TOTAL', memes.data.hits.total)
  },
}
