import { create } from "zustand";
import { persist } from "zustand/middleware";

interface User {
    id: number;
    email: string;
    full_name: string;
    nickname?: string;
    image_url?: string;
    position_id?: string;
    position_name?: string;
    committee_id?: string;
    committee_name?: string;
    division_id?: string;
    division_name?: string;
    house_name?: string;
    contact_number?: string;
    college?: string;
    program?: string;
    interests?: string;
    discord?: string;
    fb_link?: string;
    telegram?: string;
}

interface AuthState {
    user: User | null;
    isLoading: boolean;
    isAuthenticated: boolean;
    setUser: (user: User | null) => void;
    setLoading: (loading: boolean) => void;
    logout: () => void;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set) => ({
            user: null,
            isLoading: true,
            isAuthenticated: false,

            setUser: (user) =>
                set({
                    user,
                    isAuthenticated: !!user,
                    isLoading: false,
                }),

            setLoading: (isLoading) => set({ isLoading }),

            logout: () =>
                set({
                    user: null,
                    isAuthenticated: false,
                    isLoading: false,
                }),
        }),
        {
            name: "auth-storage",
            partialize: (state) => ({ user: state.user }),
        },
    ),
);
