import React, { useState, useEffect } from 'react';
import { getEventAppsForAdmin } from '../services/api';

const EventApplications = () => {
    const [eventApps, setEventApps] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchEventApps = async () => {
            try {
                const response = await getEventAppsForAdmin();
                setEventApps(response.data.event_apps);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchEventApps();
    }, []);

    if (loading) {
        return <div>Loading event applications...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className="p-4">
            <h1 className="text-2xl font-bold mb-4">Event Applications</h1>
            {eventApps.map((eventApp) => (
                <div key={eventApp.id} className="mb-6 p-4 border rounded-lg">
                    <h2 className="text-xl font-semibold mb-2">
                        Event ID: {eventApp.event_id}
                    </h2>
                    <div className="ml-4">
                        <h3 className="font-medium mb-2">Attendees:</h3>
                        <ul className="list-disc list-inside">
                            {eventApp.attendee.map((email, index) => (
                                <li key={index} className="text-gray-700">
                                    {email}
                                </li>
                            ))}
                        </ul>
                    </div>
                </div>
            ))}
        </div>
    );
};

export default EventApplications;