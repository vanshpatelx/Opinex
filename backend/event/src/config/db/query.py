class Query:
    """
    SQL queries for fetching order and event data.

    Queries:
        1. `get_event`: Fetch details of a specific event by its ID.
        2. `get_all_event`: Fetch a paginated list of running events.

    Notes:
        - Queries use positional parameters ($1, $2, etc.) for safe parameter substitution.
        - Results are ordered by timestamp in descending order.

    Last Updated: February 28, 2025
    """

    get_event = """
        SELECT id, name, details, status, settlement_time 
        FROM "Event" WHERE id = $1
    """

    get_all_event = """
        SELECT id, name, details, status, settlement_time 
        FROM "Event" WHERE status = 'RUNNING' 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2
    """
