export default function (profile?: any): globalThis.Ref<Profile | null, Profile | null> {
    const profile_ = useState<Profile | null>('profile', () => null)
    if (profile != null && JSON.stringify(profile) != JSON.stringify(profile_)) {
        profile_.value = profile
    }
    return profile_
}

interface Profile {
    id: string,
    name: string,
    email: string,
    github: string,
    linkedin: string,
    web: string,
    description: string,
    is_admin: number,
    created_at: number,
    update_at: number,
}