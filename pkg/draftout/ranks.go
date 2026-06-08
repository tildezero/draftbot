package draftout

import "context"

func (c *Client) GetRanks(ctx context.Context) ([]Rank, error) {
	var resu []Rank
	err := c.get(ctx, "/api/ranks", nil, &resu)
	if err != nil {
		return nil, err
	}
	return resu, nil
}
