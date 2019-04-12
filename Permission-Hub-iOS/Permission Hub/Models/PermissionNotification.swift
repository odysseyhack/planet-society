//
//  PermissionNotification.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

struct PermissionItem: Decodable {
    let item: String
    let fields: [String]

    private enum CodingKeys: String, CodingKey {
        case item = "Item"
        case fields = "Fields"
    }
}

struct PermissionNotification: Decodable {
    let transactionID: String
    let item: [PermissionItem]
    let reason: String
    let verification: [String]
    let date: Date
    let requesterName: String
    let requesterPublicKey: String

    private enum CodingKeys: String, CodingKey {
        case transactionID
        case item
        case reason
        case verification
        case date
        case requesterName
        case requesterPublicKey = "RequesterPublicKey"
    }

    // MARK: - Initialization

    init(
        transactionID: String,
        item: [PermissionItem],
        reason: String,
        verification: [String],
        date: Date,
        requesterName: String,
        requesterPublicKey: String) {

        self.transactionID = transactionID
        self.item = item
        self.reason = reason
        self.verification = verification
        self.date = date
        self.requesterName = requesterName
        self.requesterPublicKey = requesterPublicKey
    }

    init(from decoder: Decoder) throws {

        let container = try decoder.container(keyedBy: CodingKeys.self)
        let transactionID = try container.decode(String.self, forKey: .transactionID)
        let item = try container.decode([PermissionItem].self, forKey: .item)
        let reason = try container.decode(String.self, forKey: .reason)
        let verification = try container.decode([String].self, forKey: .verification)

        let dateString = try container.decode(String.self, forKey: .date)
        let dateFormatter = ISO8601DateFormatter()

        guard let date = dateFormatter.date(from: dateString) else {
            throw DecodingError.valueNotFound(
                Date.self,
                DecodingError.Context(
                    codingPath: [],
                    debugDescription: "Error decoding date object"))
        }

        let requesterName = try container.decode(String.self, forKey: .requesterName)
        let requesterPublicKey = try container.decode(String.self, forKey: .requesterPublicKey)

        self.init(
            transactionID: transactionID,
            item: item,
            reason: reason,
            verification: verification,
            date: date,
            requesterName: requesterName,
            requesterPublicKey: requesterPublicKey)
    }
}
