import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/designkit.dart' as ds;
import 'feed.edit.dart';
import 'feed.row.dart';
import 'api.dart' as api;

class Item extends StatelessWidget {
  final api.Feed current;
  final void Function(api.Feed? upd) onChange;
  const Item({
    super.key,
    required api.Feed this.current,
    this.onChange = ds.fnNoop,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ds.TableRow.single(
      FeedRow(current: current, onChange: onChange),
      expanded: Container(
        padding: theme.buttonTheme.padding,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Edit(
              current: current,
              onChange: (u) {
                httpx
                    .withRetry(
                      () => api.create(
                        api.FeedCreateRequest(feed: u),
                        options: [authn.request(authn.AuthzCache.meta(context))],
                      ),
                    )
                    .then((resp) {
                      onChange(resp.feed);
                    });
              },
            ),
          ],
        ),
      ),
    );
  }
}
