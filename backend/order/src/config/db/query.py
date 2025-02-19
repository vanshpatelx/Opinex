class Query:
    """
    SQL queries for fetching order data.

    Queries:
        1. `get_order_by_orderID`: Fetch a specific order by its ID.
        2. `get_order_by_eventID_pagination_admin`: Fetch paginated orders for an event (Admin).
        3. `get_order_by_eventID_pagination_user`: Fetch paginated orders for an event (User-specific).
        4. `get_order_by_userID_pagination`: Fetch paginated orders for a user.

    Notes:
        - Queries use positional parameters ($1, $2, etc.) for safe parameter substitution.
        - Results are ordered by timestamp in descending order.

    Last Updated: February 19, 2025
    """

    get_order_by_orderID = "SELECT * FROM orders WHERE id = $1"

    get_order_by_eventID_pagination_admin = """
    SELECT * FROM orders WHERE eventid = $1
    ORDER BY timestamp DESC LIMIT $2 OFFSET $3
    """

    get_order_by_eventID_pagination_user = """
    SELECT * FROM orders WHERE userid = $1 AND eventid = $2 
    ORDER BY timestamp DESC LIMIT $3 OFFSET $4
    """

    get_order_by_userID_pagination = """
    SELECT * FROM orders WHERE userid = $1
    ORDER BY timestamp DESC LIMIT $2 OFFSET $3
    """
