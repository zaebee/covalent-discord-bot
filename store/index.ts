import $http from '~/plugins/axios'
import { GetterTree, ActionTree, MutationTree } from 'vuex'

export const state = () => ({
  things: [] as string[],
  memes: [] as Object[],
  name: 'Me',
})

export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  name: state => state.name,
}

export const mutations: MutationTree<RootState> = {
  CHANGE_NAME: (state, newName: string) => (state.name = newName),
  SET_MEMES: (state, memes: Object[]) => (state.memes = memes),
}

export const actions: ActionTree<RootState, RootState> = {
  async fetchMemes({ commit }) {
    const memes = await $http.post('/memes/_search', {sort: {Timestamp: {order: 'desc'}}})
    console.log(memes.data.hits.hits)
    commit('SET_MEMES', memes.data.hits.hits)
    commit('CHANGE_NAME', 'New name')
  },
}
