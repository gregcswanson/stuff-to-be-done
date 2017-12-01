var initalState = {
    user: {},
    isLoggedIn: false,
    isConnected: true, // prevent editing when offline - only allow quick adds
    token: {},
    current: {
        todayAsInt: 0,
        isDirty: false, // use to prevent over-writing failed save requests when coming out of offline mode
        days: [
            {dayAsInt: 0, tasks: []}
        ]
    },
    history: {
        dateAsInt: 0,
        tasks: []
    },
    later: {
        count: 0,
        days: [],
    },
    trash: {
        count: 0,
        days: []
    },
    offlineChanges: [], // quick add items only
};