//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum TransactionOverviewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case description(text: String)
    case transactionItem(item: TransactionItem)
}

final class TransactionOverviewViewController: UITableViewController {

    // MARK: - Private properties

    private let transaction: TransactionNotification

    private var cells: [TransactionOverviewViewCellType] {

        var cells = [TransactionOverviewViewCellType]()

        cells.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        cells.append(.description(text: transaction.reason))

        transaction.items.forEach {
            cells.append(.transactionItem(item: $0))
        }

        return cells
    }

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(style: .plain)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        configureTableView()
    }

    // MARK: - Configuration

    private func configureTableView() {

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            UITableViewCell.self,
            forCellReuseIdentifier: String(describing: UITableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.allowsSelection = false
    }

    // MARK: - UITableViewDataSource

    override func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return cells.count
    }

    override func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch cells[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: UITableViewCell.self),
                for: indexPath)

            cell.textLabel?.font = PHFonts.regular()
            cell.textLabel?.text = text

            return cell

        case .transactionItem(let item):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "))
            cell.configure(withViewModel: viewModel)

            return cell
        }
    }
}
