import 'package:flutter/material.dart';
import 'package:retrovibed/authn.dart' as authn;
import 'package:retrovibed/designkit.dart' as ds;
import 'package:retrovibed/httpx.dart' as httpx;
import 'package:retrovibed/media/media.pb.dart';
import 'package:retrovibed/mimex.dart' as mimex;
import 'api.dart' as api;

class KnownMediaIcon extends StatelessWidget {
  final Media media;
  final double size;

  const KnownMediaIcon(this.media, {super.key, this.size = 24});

  Map<String, String>? _imageheaders(String original) {
    if (original.isEmpty) return null;
    if (!original.startsWith("https://${httpx.host()}")) return null;
    return <String, String>{"Authorization": httpx.auto_bearer_host()};
  }

  @override
  Widget build(BuildContext context) {
    final fallback = Icon(mimex.icon(media.mimetype), size: size);
    final authz = authn.AuthzCache.meta(context);

    return FutureBuilder<api.Known>(
      future: api.known.autodetect(media, options: [authn.request(authz)]),
      builder: (context, snapshot) {
        final image = snapshot.data?.image ?? "";
        if (image.isEmpty) return fallback;

        return ClipRRect(
          borderRadius: BorderRadius.circular(size / 4),
          child:
              ds.Image.precache(
                context,
                image,
                headers: _imageheaders(image),
                width: size,
                height: size,
                fit: BoxFit.cover,
                missing: fallback,
              ) ??
              fallback,
        );
      },
    );
  }
}
