//
//  TransactionNotification.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import Foundation

struct TransactionItem: Decodable {
    let item: String
    let fields: [String]

    private enum CodingKeys: String, CodingKey {
        case item = "Item"
        case fields = "Fields"
    }
}

struct TransactionNotification: Decodable {
    let transactionID: String
    let items: [TransactionItem]
    let title: String
    let description: String
    let verification: [String]
    let date: Date
    let requesterName: String
    let requesterPublicKey: String
    let analysis: [String]

    private enum CodingKeys: String, CodingKey {
        case transactionID
        case items = "item"
        case title
        case description
        case verification
        case date
        case requesterName
        case requesterPublicKey = "RequesterPublicKey"
        case analysis
    }
}
