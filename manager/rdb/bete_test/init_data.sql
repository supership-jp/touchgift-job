INSERT INTO store (
            organization_code,
            store_id,
            name,
            post_code,
            province_code,
            address
        ) VALUES ('org_1', 'store1', 'store_name1', '0000000', '00', '住所');
INSERT INTO store_group (
			name,
		 	organization_code,
			last_updated_by
        ) VALUES ('group1', 'org_1', 1);
INSERT INTO gimmick (
            name,
            img_url,
            organization_code,
            status,
						code,
            xid,
            last_updated_by
        ) VALUES ('gimmick1', 'http://example.com', 'org_1', '2', 'code1', 'xid', 1);
INSERT INTO campaign (
            organization_code,
            status,
            name,
            start_at,
            end_at,
            created_at,
            updated_at,
            last_updated_by,
            store_group_id,
						daily_coupon_limit_per_user
        ) VALUES ('org_1', 'configured', 'campaign1', '2024-08-27 16:30:13', '2024-08-27 17:46:13', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, 1, 1, 1);
INSERT INTO coupon (
            name,
            organization_code,
            status,
            code,
            img_url,
            xid
        ) VALUES ('coupon1', 'org_1', '2', 'code1', 'http://example.com', 'xid');
INSERT INTO campaign_coupon (
            campaign_id,
            coupon_id,
            delivery_rate
        ) VALUES (1, 1, 100);
INSERT INTO campaign_gimmick (
						campaign_id,
						gimmick_id
				) VALUES (1, 1);
INSERT INTO video (
            video_url,
            endcard_url,
            video_xid,
            endcard_xid,
            height,
            width,
            extension,
            endcard_height,
            endcard_width,
		    endcard_extension,
			last_updated_by,
		    duration,
		    endcard_link
        ) VALUES ('video_url', 'endcard_url', 'video_xid', 'endcard_xid', 100, 100, 'mp4', 100, 100, 'png', 1, 100, 'http://example.com');
INSERT INTO creative (
            organization_code,
            status,
						name,
            click_url,
            last_updated_by,
            video_id
        ) VALUES ('org_1','2','creative1','http://example.com',1,1);
INSERT INTO campaign_creative (
			campaign_id,
		    creative_id,
		    delivery_rate,
		    skip_offset
        ) VALUES (1,1,100, 1);
INSERT INTO touch_point (
				  organization_code,
				  point_unique_id,
				  print_management_id,
				  store_id,
				  type,
				  name,
				  last_updated_by
				) VALUES (
				    'org_1', 'point_unique_id', 'print_management_id', 1, 'QR', 'name', 1
				);
INSERT INTO store_map(
		     store_group_id,
			 store_id
		) VALUES (1, 1);