"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useAuthStore } from "@/lib/auth-store";

export function useAuth() {
    const { user, setUser, isAuthenticated, isLoading, logout } =
        useAuthStore();
    const queryClient = useQueryClient();

    const { data: member, isLoading: isMemberLoading } = useQuery({
        queryKey: ["member", "me"],
        queryFn: () => api.getMe(),
        enabled: isAuthenticated,
        staleTime: 5 * 60 * 1000, // 5 minutes
    });

    React.useEffect(() => {
        if (member) {
            setUser(member);
        }
    }, [member, setUser]);

    const loginMutation = useMutation({
        mutationFn: () => Promise.resolve(api.login()),
    });

    const logoutMutation = useMutation({
        mutationFn: () => api.logout(),
        onSuccess: () => {
            queryClient.clear();
            logout();
        },
    });

    const updateProfileMutation = useMutation({
        mutationFn: (data: Parameters<typeof api.updateMe>[0]) =>
            api.updateMe(data),
        onSuccess: (updatedUser) => {
            setUser(updatedUser);
            queryClient.invalidateQueries({ queryKey: ["member", "me"] });
        },
    });

    return {
        user,
        isLoading: isLoading || isMemberLoading,
        isAuthenticated,
        login: loginMutation.mutateAsync,
        logout: logoutMutation.mutateAsync,
        updateProfile: updateProfileMutation.mutateAsync,
    };
}

import React from "react";
